const fs = require('fs');
const path = require('path');
const { spawnSync } = require('child_process');

const modules = [
  'abs',
  'acos',
  'acosh',
  'asin',
  'asinh',
  'atan',
  'atan2',
  'atanh',
  'cbrt',
  'ceil',
  'clamp',
  'clone',
  'cmp',
  'config',
  'cos',
  'cosh',
  'Decimal',
  'div',
  'divToInt',
  'dpSd',
  'exp',
  'floor',
  'hypot',
  'immutability',
  'intPow',
  'isFiniteEtc',
  'ln',
  'log',
  'log10',
  'log2',
  'minAndMax',
  'minus',
  'mod',
  'neg',
  'plus',
  'pow',
  'random',
  'round',
  'sign',
  'sin',
  'sinh',
  'sqrt',
  'sum',
  'tan',
  'tanh',
  'times',
  'toBinary',
  'toDP',
  'toExponential',
  'toFixed',
  'toFraction',
  'toHex',
  'toNearest',
  'toNumber',
  'toOctal',
  'toPrecision',
  'toSD',
  'toString',
  'trunc',
  'valueOf'
];

const csvPath = path.join(__dirname, '../../../audit-reports/FAILURE_DATABASE.csv');
const csvHeaders = [
  'Module',
  'Assertion number',
  'Input(s)',
  'Precision',
  'Rounding mode',
  'Expected output',
  'Actual output',
  'JS stack',
  'Go operation',
  'First differing function',
  'First differing source line',
  'Assigned cluster',
  'Status'
];

const records = [];

function escapeCSV(val) {
  if (val === undefined || val === null) return '';
  let str = String(val).replace(/\r/g, '').trim();
  if (str.includes(',') || str.includes('"') || str.includes('\n')) {
    str = '"' + str.replace(/"/g, '""') + '"';
  }
  return str;
}

function getSourceLine(filePath, lineNum) {
  try {
    if (!fs.existsSync(filePath)) return '';
    const content = fs.readFileSync(filePath, 'utf8');
    const lines = content.split('\n');
    if (lineNum > 0 && lineNum <= lines.length) {
      return lines[lineNum - 1].trim();
    }
  } catch (e) {}
  return '';
}

function assignCluster(moduleName, errorText) {
  const txt = (errorText || '').toLowerCase();
  if (txt.includes('is not a function') || txt.includes('undefined') || txt.includes('not defined')) {
    return 'Cluster 1 (Missing APIs)';
  }
  if (moduleName === 'cbrt') {
    return 'Cluster 2 (cbrt)';
  }
  if (moduleName === 'cmp') {
    return 'Cluster 3 (cmp)';
  }
  if (moduleName === 'cos' || moduleName === 'sin' || moduleName === 'tan') {
    return 'Cluster 4 (Trig reduction)';
  }
  if (moduleName === 'config') {
    return 'Cluster 5 (config)';
  }
  if (['cosh', 'sinh', 'asinh', 'atanh', 'acosh'].includes(moduleName)) {
    return 'Cluster 6 (Hyperbolic)';
  }
  if (moduleName === 'exp') {
    return 'Cluster 7 (exp)';
  }
  if (['asin', 'acos', 'atan'].includes(moduleName)) {
    return 'Cluster 8 (Inverse Trig)';
  }
  if (txt.includes('nan') || txt.includes('infinity') || txt.includes('-0')) {
    return 'Cluster 9 (Special values)';
  }
  if (['toFixed', 'toExponential', 'toPrecision', 'toSignificantDigits', 'toFraction', 'toNumber', 'toBinary', 'toOctal', 'toHex', 'toNearest', 'toPrecision', 'toSD', 'hypot', 'log', 'log2', 'log10', 'sign', 'sum', 'minAndMax', 'random'].includes(moduleName)) {
    return 'Cluster 10 (Formatting/Query)';
  }
  return 'Other Math/Logic';
}

console.log(`Scanning ${modules.length} modules for failures...`);

for (const moduleName of modules) {
  console.log(`Running module: ${moduleName}`);
  
  // Prepare code execution block
  const code = `
    const Decimal = require('./bridge');
    global.Decimal = Decimal;
    Decimal.config({
      precision: 20,
      rounding: 4,
      toExpNeg: -7,
      toExpPos: 21,
      minE: -9e15,
      maxE: 9e15
    });
    require('./setup');
    try {
      require('./modules/${moduleName}');
    } catch(err) {
      console.error('CRASH_ERROR:' + err.stack);
      process.exit(1);
    }
  `;

  const res = spawnSync('node', ['-e', code], {
    cwd: __dirname,
    encoding: 'utf8',
    timeout: 15000 // 15s timeout
  });

  if (res.error && res.error.code === 'ETIMEDOUT') {
    records.push({
      Module: moduleName,
      'Assertion number': 'N/A',
      'Input(s)': 'Timeout after 15s',
      Precision: 'N/A',
      'Rounding mode': 'N/A',
      'Expected output': 'N/A',
      'Actual output': 'Timeout',
      'JS stack': 'N/A',
      'Go operation': 'N/A',
      'First differing function': 'N/A',
      'First differing source line': 'N/A',
      'Assigned cluster': assignCluster(moduleName, 'timeout'),
      Status: 'Open'
    });
    continue;
  }

  const stdout = res.stdout || '';
  const stderr = res.stderr || '';

  // Check for crash error in stderr/stdout
  if (stderr.includes('CRASH_ERROR:') || stdout.includes('CRASH_ERROR:') || res.status !== 0) {
    const combined = stdout + '\n' + stderr;
    const crashIdx = combined.indexOf('CRASH_ERROR:');
    let stack = combined;
    if (crashIdx !== -1) {
      stack = combined.substring(crashIdx + 12);
    }
    const firstLine = stack.split('\n')[0] || 'Unknown error';
    
    // Find if there is a specific line in modules directory that threw
    let lineNum = 'N/A';
    let sourceLine = 'N/A';
    let fileMatch = stack.match(/modules\/([a-zA-Z0-9_\.]+):(\d+)/);
    if (fileMatch) {
      const fileName = fileMatch[1];
      const lineNo = parseInt(fileMatch[2], 10);
      lineNum = lineNo;
      sourceLine = getSourceLine(path.join(__dirname, 'modules', fileName), lineNo);
    }

    records.push({
      Module: moduleName,
      'Assertion number': 'N/A',
      'Input(s)': sourceLine !== 'N/A' ? sourceLine : 'Module Execution Crash',
      Precision: 'N/A',
      'Rounding mode': 'N/A',
      'Expected output': 'N/A',
      'Actual output': firstLine,
      'JS stack': stack,
      'Go operation': 'N/A',
      'First differing function': 'N/A',
      'First differing source line': lineNum !== 'N/A' ? `line ${lineNum}` : 'N/A',
      'Assigned cluster': assignCluster(moduleName, firstLine),
      Status: 'Open'
    });
    continue;
  }

  // Parse assertion failures from stdout
  const lines = stdout.split('\n');
  let i = 0;
  while (i < lines.length) {
    const line = lines[i];
    if (line.includes('failed:')) {
      const match = line.match(/Test number (\d+) failed: (\w+)/);
      if (match) {
        const testNum = match[1];
        const assertType = match[2];
        let expected = '';
        let actual = '';
        let stackLines = [];
        
        i++;
        while (i < lines.length && !lines[i].includes('failed:') && !lines[i].startsWith(' Testing') && !lines[i].includes('tests passed in')) {
          const l = lines[i];
          if (l.trim().startsWith('Expected:')) {
            expected = l.substring(l.indexOf('Expected:') + 9).trim();
          } else if (l.trim().startsWith('Actual:')) {
            actual = l.substring(l.indexOf('Actual:') + 7).trim();
          } else if (l.trim().startsWith('x:')) {
            expected = 'x: ' + l.substring(l.indexOf('x:') + 2).trim();
          } else if (l.trim().startsWith('y:')) {
            actual = 'y: ' + l.substring(l.indexOf('y:') + 2).trim();
          } else if (l.trim().startsWith('Expected digits:')) {
            expected += 'digits: ' + l.substring(l.indexOf('Expected digits:') + 16).trim() + '; ';
          } else if (l.trim().startsWith('Expected exponent:')) {
            expected += 'exp: ' + l.substring(l.indexOf('Expected exponent:') + 18).trim() + '; ';
          } else if (l.trim().startsWith('Expected sign:')) {
            expected += 'sign: ' + l.substring(l.indexOf('Expected sign:') + 14).trim();
          } else if (l.trim().startsWith('Actual digits:')) {
            actual += 'digits: ' + l.substring(l.indexOf('Actual digits:') + 14).trim() + '; ';
          } else if (l.trim().startsWith('Actual exponent:')) {
            actual += 'exp: ' + l.substring(l.indexOf('Actual exponent:') + 16).trim() + '; ';
          } else if (l.trim().startsWith('Actual sign:')) {
            actual += 'sign: ' + l.substring(l.indexOf('Actual sign:') + 12).trim();
          } else if (l.trim().startsWith('Stack:')) {
            // skip
          } else if (l.trim() !== '') {
            stackLines.push(l);
          }
          i++;
        }
        
        const stackStr = stackLines.join('\n');
        
        // Extract line number in modules
        let lineNum = 'N/A';
        let sourceLine = 'N/A';
        const fileMatch = stackStr.match(/modules\/([a-zA-Z0-9_\.]+):(\d+)/);
        if (fileMatch) {
          const fileName = fileMatch[1];
          const lineNo = parseInt(fileMatch[2], 10);
          lineNum = lineNo;
          sourceLine = getSourceLine(path.join(__dirname, 'modules', fileName), lineNo);
        }
        
        records.push({
          Module: moduleName,
          'Assertion number': testNum,
          'Input(s)': sourceLine !== 'N/A' ? sourceLine : 'N/A',
          Precision: 'N/A',
          'Rounding mode': 'N/A',
          'Expected output': expected,
          'Actual output': actual,
          'JS stack': stackStr,
          'Go operation': moduleName,
          'First differing function': 'N/A',
          'First differing source line': lineNum !== 'N/A' ? `line ${lineNum}` : 'N/A',
          'Assigned cluster': assignCluster(moduleName, expected + ' ' + actual),
          Status: 'Open'
        });
        
        i--;
      }
    }
    i++;
  }
}

// Format as CSV
let csvContent = csvHeaders.join(',') + '\n';
for (const rec of records) {
  const row = csvHeaders.map(h => escapeCSV(rec[h]));
  csvContent += row.join(',') + '\n';
}

fs.writeFileSync(csvPath, csvContent, 'utf8');
console.log(`Successfully generated ${records.length} records in FAILURE_DATABASE.csv`);
