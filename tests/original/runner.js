const fs = require('fs');
const path = require('path');

const modulesDir = path.join(__dirname, 'modules');
const files = fs.readdirSync(modulesDir).filter(f => f.endsWith('.js'));

console.log(`Executing ${files.length} original decimal.js test modules against Go decimal binary bridge...\n`);

let totalPassed = 0;
let totalFailed = 0;

for (const file of files) {
  const filePath = path.join(modulesDir, file);
  try {
    require(filePath);
  } catch (err) {
    console.error(`Error executing ${file}:`, err.message);
    totalFailed++;
  }
}

console.log(`\nOriginal Test Suite Harness Run Complete.`);
