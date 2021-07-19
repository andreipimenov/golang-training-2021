# Homework 02

Given mathematical expressions as a string, e.g
```
20/2-(2+2*3)=
```

Supported operators: `+`, `-`, `*`, `/`, `(`, `)`.

Implement `Calc` interface which supports those operators in accordance with their priorities.

```
// Division
20/2 = 10

// Expression in parentheses
2*3 = 6
2+6 = 8

// Subtraction
10-8 = 2
```

So, `20/2-(2+2*3)=` should return `2`

Keep in mind that it should work with other examples.

The following interface should be implemented

```go
type Calc interface {
    Calculate(expression string) float64
}
```

**Additionally:**  
- Add validation for your app to check if no symbols other than digits and allowed operators are passed as an argument to `Calculate(string)` method
- Use structs and methods, try to implement the task using tree-based approach and stack-based approach
