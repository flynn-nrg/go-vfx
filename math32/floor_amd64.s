//go:build amd64

#include "textflag.h"

// func floor(x float32) float32
TEXT ·floor(SB),NOSPLIT,$0-8
	// ABI wrapper stores input at 8(SP) and expects output at 16(SP)
	MOVSS	8(SP), X0        // Load x from stack
	ROUNDSS	$1, X0, X0       // Round toward -Inf (floor), mode = 0x01
	MOVSS	X0, 16(SP)       // Store result
	RET

// func ceil(x float32) float32
TEXT ·ceil(SB),NOSPLIT,$0-8
	// ABI wrapper stores input at 8(SP) and expects output at 16(SP)
	MOVSS	8(SP), X0        // Load x from stack
	ROUNDSS	$2, X0, X0       // Round toward +Inf (ceil), mode = 0x02
	MOVSS	X0, 16(SP)       // Store result
	RET

// func round(x float32) float32
// ROUNDSS doesn't have "ties away from zero" mode, so we use a simpler approach
TEXT ·round(SB),NOSPLIT,$0-8
	MOVSS	8(SP), X0        // Load x from stack
	
	// Use ROUNDSS mode 0 (round to nearest, ties to even) and adjust
	// We'll use a different approach: check sign, add 0.5, truncate
	
	// Get absolute value and sign
	MOVSS	X0, X1           // Copy x
	MOVL	$0x7FFFFFFF, AX  // Mask for abs (clear sign bit)
	MOVD	AX, X2
	ANDPS	X2, X1           // X1 = |x|
	
	// Add 0.5
	MOVL	$0x3F000000, AX  // 0.5
	MOVD	AX, X2
	ADDSS	X2, X1           // X1 = |x| + 0.5
	
	// Truncate
	ROUNDSS	$3, X1, X1       // Round toward zero
	
	// Copy sign from original x
	MOVL	$0x80000000, AX  // Sign bit mask
	MOVD	AX, X2
	ANDPS	X2, X0           // X0 = sign bit of x
	ORPS	X0, X1           // X1 = result with correct sign
	
	MOVSS	X1, 16(SP)       // Store result
	RET

