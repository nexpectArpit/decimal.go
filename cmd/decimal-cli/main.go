package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	decimal "our-projectInGO/src"
)

type Request struct {
	Op        string        `json:"op"`
	Precision int           `json:"precision"`
	Rounding  int           `json:"rounding"`
	ToExpNeg  *int64        `json:"toExpNeg,omitempty"`
	ToExpPos  *int64        `json:"toExpPos,omitempty"`
	MinE      *int64        `json:"minE,omitempty"`
	MaxE      *int64        `json:"maxE,omitempty"`
	Args      []interface{} `json:"args"`
}

type Response struct {
	Sign     int8    `json:"s"`
	Exponent int64   `json:"e"`
	Digits   []int32 `json:"d"`
	String   string  `json:"str"`
	CmpRes   int     `json:"cmpRes"`
	BoolRes  bool    `json:"boolRes"`
	IsCmp    bool    `json:"isCmp"`
	IsBool   bool    `json:"isBool"`
	Error    string  `json:"error,omitempty"`
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--rpc" {
		runRPC()
		return
	}

	fmt.Println("decimal-cli v1.0.0 (Go implementation of decimal.js)")
}

func runRPC() {
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)

	for {
		var req Request
		if err := dec.Decode(&req); err != nil {
			if err == io.EOF {
				break
			}
			enc.Encode(Response{Error: err.Error()})
			continue
		}

		ctx := decimal.DefaultContext()
		if req.Precision > 0 {
			ctx.Config(decimal.WithPrecision(req.Precision))
		}
		if req.Rounding >= 0 {
			ctx.Config(decimal.WithRounding(decimal.RoundingMode(req.Rounding)))
		}
		if req.ToExpNeg != nil {
			ctx.Config(decimal.WithToExpNeg(*req.ToExpNeg))
		}
		if req.ToExpPos != nil {
			ctx.Config(decimal.WithToExpPos(*req.ToExpPos))
		}
		if req.MinE != nil {
			ctx.Config(decimal.WithMinE(*req.MinE))
		}
		if req.MaxE != nil {
			ctx.Config(decimal.WithMaxE(*req.MaxE))
		}

		resp := handleOp(ctx, req)
		enc.Encode(resp)
	}
}

func handleOp(ctx *decimal.Context, req Request) Response {
	if len(req.Args) == 0 {
		return Response{Error: "no arguments provided"}
	}

	a, err := ctx.New(req.Args[0])
	if err != nil {
		return Response{Error: err.Error()}
	}

	switch req.Op {
	case "new":
		return makeResp(ctx, a)

	case "cmp", "comparedTo":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return Response{IsCmp: true, CmpRes: a.Cmp(b)}

	case "eq", "equals":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return Response{IsBool: true, BoolRes: a.Eq(b)}

	case "gt", "greaterThan":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return Response{IsBool: true, BoolRes: a.Gt(b)}

	case "gte", "greaterThanOrEqualTo":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return Response{IsBool: true, BoolRes: a.Gte(b)}

	case "lt", "lessThan":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return Response{IsBool: true, BoolRes: a.Lt(b)}

	case "lte", "lessThanOrEqualTo":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return Response{IsBool: true, BoolRes: a.Lte(b)}

	case "add", "plus":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return makeResp(ctx, ctx.Add(a, b))

	case "sub", "minus":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return makeResp(ctx, ctx.Sub(a, b))

	case "mul", "times":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return makeResp(ctx, ctx.Mul(a, b))

	case "div", "dividedBy":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return makeResp(ctx, ctx.Div(a, b))

	case "mod", "modulo":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return makeResp(ctx, ctx.Mod(a, b))

	case "sqrt":
		return makeResp(ctx, ctx.Sqrt(a))

	case "cbrt":
		return makeResp(ctx, ctx.Cbrt(a))

	case "pow":
		if len(req.Args) < 2 {
			return Response{Error: "missing operand"}
		}
		b, err := ctx.New(req.Args[1])
		if err != nil {
			return Response{Error: err.Error()}
		}
		return makeResp(ctx, ctx.Pow(a, b))

	case "ln":
		return makeResp(ctx, ctx.Ln(a))

	case "exp":
		return makeResp(ctx, ctx.Exp(a))

	case "sin":
		return makeResp(ctx, ctx.Sin(a))

	case "cos":
		return makeResp(ctx, ctx.Cos(a))

	case "tan":
		return makeResp(ctx, ctx.Tan(a))

	case "asin":
		return makeResp(ctx, ctx.Asin(a))

	case "acos":
		return makeResp(ctx, ctx.Acos(a))

	case "atan":
		return makeResp(ctx, ctx.Atan(a))

	case "sinh":
		return makeResp(ctx, ctx.Sinh(a))

	case "cosh":
		return makeResp(ctx, ctx.Cosh(a))

	case "tanh":
		return makeResp(ctx, ctx.Tanh(a))

	case "asinh":
		return makeResp(ctx, ctx.Asinh(a))

	case "acosh":
		return makeResp(ctx, ctx.Acosh(a))

	case "atanh":
		return makeResp(ctx, ctx.Atanh(a))

	case "abs":
		return makeResp(ctx, a.Abs())

	case "neg":
		return makeResp(ctx, a.Neg())

	case "trunc", "truncated":
		return makeResp(ctx, ctx.Trunc(a))

	case "floor":
		return makeResp(ctx, ctx.Floor(a))

	case "ceil":
		return makeResp(ctx, ctx.Ceil(a))

	case "round":
		return makeResp(ctx, ctx.Round(a))

	default:
		return makeResp(ctx, a)
	}
}

func makeResp(ctx *decimal.Context, d *decimal.Decimal) Response {
	if d == nil {
		return Response{Sign: 0, Exponent: 0, Digits: nil, String: "NaN"}
	}
	return Response{
		Sign:     d.Sign(),
		Exponent: d.Exponent(),
		Digits:   d.Coefficients(),
		String:   ctx.String(d),
	}
}
