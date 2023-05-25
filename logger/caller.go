package logger

import (
	"runtime"
	"strings"
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

// Positions in the call stack when tracing to report the calling method
var minimumCallerDepth int

//var moduleName string

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

// getCaller retrieves the name of the first non-logrus calling function
func getCallers() *runtime.Frame {
	// cache this package's fully-qualified name
	getPackageNameOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)
		var de int
		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()

			if strings.Contains(funcName, "getCallers") {
				packageName = getPackageName(funcName)
				de = i
				//moduleName = getModuleName(packageName)
				break
			}
		}

		//minimumCallerDepth = knownLogrusFrames
		minimumCallerDepth = de
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		// If the caller isn't part of this package, we're done
		if pkg != packageName && pkg != "github.com/sirupsen/logrus" {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

//func getModuleName(pkg string) string {
//	if strings.HasPrefix(pkg, "github.com") {
//		return pkg
//	}
//	firstSlash := strings.Index(pkg, "/")
//	if firstSlash == -1 {
//		return pkg
//	}
//	return pkg[:firstSlash]
//}

//
//func getFuncCaller(file string) string {
//	firstSlash := strings.Index(file, "molo")
//	fmt.Println("moduleName", moduleName)
//	fmt.Println("fi", firstSlash)
//	if firstSlash == -1 {
//		return file
//	}
//	return file[firstSlash:]
//}
