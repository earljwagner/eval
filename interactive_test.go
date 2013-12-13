package interactive

// Utilities for other tests live here

import (
	"strings"
	"testing"
	"reflect"

	"go/parser"
)

func expectVoid(t *testing.T, expr string, env *Env) {
	expectResults(t, expr, env, &[]interface{}{})
}

func expectNil(t *testing.T, expr string, env *Env) {
	expectResults(t, expr, env, nil)
}

func expectResult(t *testing.T, expr string, env *Env, expected interface{}) {
	expect2 := []interface{}{expected}
	expectResults(t, expr, env, &expect2)
}

func expectResults(t *testing.T, expr string, env *Env, expected *[]interface{}) {
	ctx := &Ctx{expr}
	if e, err := parser.ParseExpr(expr); err != nil {
		t.Fatalf("Failed to parse expression '%s' (%v)", expr, err)
	} else if aexpr, errs := checkExpr(ctx, e, env); errs != nil {
		t.Fatalf("Failed to check expression '%s' (%v)", expr, errs)
	} else if results, _, err := EvalExpr(ctx, aexpr, env); err != nil {
		t.Fatalf("Error evaluating expression '%s' (%v)", expr, err)
	} else {
		if nil == results {
			if expected != nil {
				t.Fatalf("Expression '%s' is nil but expected '%+v'", expr, *expected)
			}
			return
		} else if expected == nil {
			t.Fatalf("Expression '%s'expected is '%+v', expected to be nil", expr, *results)
		}
		resultsi := make([]interface{}, len(*results))
		for i, result := range *results {
			resultsi[i] = result.Interface()
		}
		if !reflect.DeepEqual(resultsi, *expected) {
			t.Fatalf("Expression '%s' yielded '%+v', expected '%+v'", expr, resultsi, *expected)
		}
	}
}

func expectError(t *testing.T, expr string, env *Env, errorString string) {
	ctx := &Ctx{expr}
	if e, err := parser.ParseExpr(expr); err != nil {
		t.Fatalf("Failed to parse expression '%s' (%v)", expr, err)
	} else if aexpr, errs := checkExpr(ctx, e, env); errs != nil {
		// TODO handle check errors
		panic("No tests should fail here (yet)")
	} else if _, _, err := EvalExpr(ctx, aexpr, env); err == nil {
		t.Fatalf("Expected expression '%s' to fail", expr)
	// Catch dogdy error messages which panic on format
	} else if err.Error() != errorString {
		t.Fatalf("Error `%s` != Expected `%s`", err.Error(), errorString)
	}
}

// deprecated, use expectError
func expectFail(t *testing.T, expr string, env *Env) {
	ctx := &Ctx{expr}
	if e, err := parser.ParseExpr(expr); err != nil {
		t.Fatalf("Failed to parse expression '%s' (%v)", expr, err)
	} else if aexpr, errs := checkExpr(ctx, e, env); errs != nil {
		// TODO handle check errors
		panic("No tests should fail here (yet)")
	} else if _, _, err := EvalExpr(ctx, aexpr, env); err == nil {
		t.Fatalf("Expected expression '%s' to fail", expr)
	// Catch dogdy error messages which panic on format
	} else if strings.Index(err.Error(), "(PANIC=") != -1 {
		t.Fatalf("Expression '%s' failed as expected but error message panicked (%v)", expr, err)
	}
}

func expectConst(t *testing.T, expr string, env *Env, expected interface{}, expectedType reflect.Type) {
	ctx := &Ctx{expr}
	if e, err := parser.ParseExpr(expr); err != nil {
		t.Fatalf("Failed to parse expression '%s' (%v)", expr, err)
	} else if aexpr, errs := checkExpr(ctx, e, env); errs != nil {
		t.Fatalf("Failed to check expression '%s' (%v)", expr, errs)
	} else if !aexpr.IsConst() {
		t.Fatalf("Expression '%s' did not yield a const node(%+v)", expr, aexpr)
	} else if expectedBigComplex, ok := expected.(*BigComplex); ok {
		if actual, ok2 := aexpr.Const().Interface().(*BigComplex); !ok2 {
			t.Fatalf("Expression '%s' yielded '%v', expected '%v'", expr, aexpr.Const(), expected)
		} else if !actual.Equals(expectedBigComplex) {
			t.Fatalf("Expression '%s' yielded '%v', expected '%v'", expr, actual, expected)
		} else if len(aexpr.KnownType()) == 0 {
			t.Fatalf("Expression '%s' expected to have type '%v'", expr, expectedType)
		} else if actual := aexpr.KnownType()[0]; !reflect.DeepEqual(actual, expectedType) {
			t.Fatalf("Expression '%s' has type '%v', expected '%v'", expr, actual, expectedType)
		}
	} else {
		if actual := aexpr.Const().Interface(); !reflect.DeepEqual(actual, expected) {
			t.Fatalf("Expression '%s' yielded '%+v', expected '%+v'", expr, actual, expected)
		} else if len(aexpr.KnownType()) == 0 {
			t.Fatalf("Expression '%s' expected to have type '%v'", expr, expectedType)
		} else if actual := aexpr.KnownType()[0]; !reflect.DeepEqual(actual, expectedType) {
			t.Fatalf("Expression '%s' has type '%v', expected '%v'", expr, t, expectedType)
		}
	}
}

func makeEnv() *Env {
	return &Env {
		Vars: make(map[string] reflect.Value),
		Consts: make(map[string] reflect.Value),
		Funcs: make(map[string] reflect.Value),
		Types: make(map[string] reflect.Type),
		Pkgs: make(map[string] Pkg),
	}
}
