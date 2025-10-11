# Float32 math package

## Introduction
This package provides a float32 math package as a drop-in replacement for Golang's `math` to be used by the [Izpi](https://github.com/flynn-nrg/izpi) renderer.

## Implemented functions

### Core Arithmetic
* Abs - Absolute value
* Max - Maximum of two values (hardware accelerated)
* Min - Minimum of two values (hardware accelerated)
* Sqrt - Square root (hardware accelerated: ARM64 FSQRTS, AMD64 SQRTSS)

### Trigonometric Functions
* Sin - Sine (optimized polynomial)
* Cos - Cosine (optimized polynomial)
* Tan - Tangent (optimized polynomial)

### Inverse Trigonometric Functions
* Asin - Arcsine
* Acos - Arccosine
* Atan - Arctangent
* Atan2 - Two-argument arctangent

### Exponential and Logarithmic
* Exp - Exponential (e^x)
* Log - Natural logarithm 

### Rounding Functions
* Floor - Round down (hardware accelerated: ARM64 FRINTMS, AMD64 ROUNDSS)
* Ceil - Round up (hardware accelerated: ARM64 FRINTPS, AMD64 ROUNDSS)
* Round - Round to nearest (hardware accelerated: ARM64 FRINTAS)

### Utility Functions
* IsNaN - Check for Not-a-Number
* IsInf - Check for infinity
* Signbit - Check sign bit

