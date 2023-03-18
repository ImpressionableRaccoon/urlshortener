/*
Multichecker staticlint нужен для статического анализа кода
и поиска возможных ошибок

# Как запускать

	staticlint [path ...]

# Анализаторы, которые используются по-умолчанию

Стандартные из golang.org/x/tools/go/analysis/passes:
  - asmdecl reports mismatches between assembly files and Go declarations;
  - assign detects useless assignments;
  - atomic checks for common mistakes using the sync/atomic package;
  - atomicalign checks for non-64-bit-aligned arguments to sync/atomic functions;
  - bools detects common mistakes involving boolean operators;
  - buildssa constructs the SSA representation of an error-free package
    and returns the set of all functions within it;
  - buildtag checks build tags;
  - cgocall detects some violations of the cgo pointer passing rules;
  - composite checks for unkeyed composite literals;
  - copylock checks for locks erroneously passed by value;
  - ctrlflow provides a syntactic control-flow graph (CFG) for the body of a function;
  - deepequalerrors checks for the use of reflect.DeepEqual with error values;
  - directive checks known Go toolchain directives;
  - errorsas checks that the second argument to errors.As is a pointer to a type implementing error;
  - fieldalignment detects structs that would use less memory if their fields were sorted;
  - findcall serves as a trivial example and test of the Analysis API;
  - framepointer reports assembly code that clobbers the frame pointer before saving it;
  - httpresponse checks for mistakes using HTTP responses;
  - ifaceassert flags impossible interface-interface type assertions;
  - inspect provides an AST inspector (golang.org/x/tools/go/ast/inspector.Inspector)
    for the syntax trees of a package;
  - loopclosure checks for references to enclosing loop variables from within nested functions;
  - lostcancel checks for failure to call a context cancellation function;
  - nilfunc checks for useless comparisons against nil;
  - nilness inspects the control-flow graph of an SSA function and reports errors
    such as nil pointer dereferences and degenerate nil pointer comparisons;
  - pkgfact is a demonstration and test of the package fact mechanism;
  - printf checks consistency of Printf format strings and arguments;
  - reflectvaluecompare checks for accidentally using ==
    or reflect.DeepEqual to compare reflect.Value values;
  - shadow checks for shadowed variables;
  - shift checks for shifts that exceed the width of an integer;
  - sigchanyzer detects misuse of unbuffered signal as argument to signal.Notify;
  - sortslice checks for calls to sort.Slice that do not use a slice type as first argument;
  - stdmethods checks for misspellings in the signatures of methods similar to well-known interfaces;
  - stringintconv flags type conversions from integers to strings;
  - structtag checks struct field tags are well formed;
  - testinggoroutine report calls to (*testing.T).Fatal from goroutines started by a test;
  - tests checks for common mistaken usages of tests and examples;
  - timeformat checks for the use of time.Format or time.Parse calls with a bad format;
  - unmarshal checks for passing non-pointer or non-interface types to unmarshal and decode functions;
  - unreachable checks for unreachable code;
  - unsafeptr checks for invalid conversions of uintptr to unsafe.Pointer;
  - unusedresult checks for unused results of calls to certain pure functions;
  - unusedwrite checks for unused writes to the elements of a struct or array object;
  - usesgenerics checks for usage of generic features added in Go 1.18.

staticcheck.io:
  - анализаторы класса SA;
  - quickfix для рефакторинга кода;
  - simple для упрощения кода;
  - stylecheck для соблюдения правил стиля.

go-critic для проверок, которые почему-то отсутствуют в других линтерах.

usestdlibvars для поиска возможности использования встроенных переменных вместо магических чисел.

osexit, который будет ругаться на прямой вызов os.Exit() в функции main пакета main.
*/
package main
