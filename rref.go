package reedsolomon

import "fmt"

// AddRowRREF adds a new row to a matrix already in RREF form, and reduces the matrix to RREF again.
func AddRowRREF(m Matrix, row []byte) (Matrix, error) {
	// Find leading entry in the row to be added to the Matrix already in RREF form.
	rowLeadingEntryIndex, err := findLeadingEntryIndex(row)
	if err != nil {
		return m, fmt.Errorf("failed to find leading entry in new row: %w", err)
	}

	var (
		scale          byte
		insertionIndex int
	)
	// TODO: It might be faster to start from the bottom and work up.
	// Ref: https://github.com/getoptimum/reedsolomon/issues/3
	for insertionIndex = 0; insertionIndex < len(m); insertionIndex++ {
		// Find leading entry in current row
		currentRowLeadingEntryIndex, err := findLeadingEntryIndex(m[insertionIndex])
		if err != nil {
			return m, fmt.Errorf("[BUG]: failed to find leading entry in current row: %w", err)
		}

		if rowLeadingEntryIndex < currentRowLeadingEntryIndex {
			// This is a good spot for us to insert the new rows
			break
		} else if rowLeadingEntryIndex == currentRowLeadingEntryIndex {
			// Scale and subtract the current row from the new row, so as to zero out the new row's leading entry.
			scale = row[rowLeadingEntryIndex]
			galMulSliceXor(scale, m[insertionIndex], row, &defaultOptions)

			rowLeadingEntryIndex, err = findLeadingEntryIndex(row)
			if err != nil {
				return m, fmt.Errorf("failed to find leading entry in new row: %w", err)
			}
		}
	}

	// Scale rrefRow to start with 1
	galMulSlice(invTable[row[rowLeadingEntryIndex]], row, row, &defaultOptions)

	// Add the new rows in between insertionIndex - 1 and insertionIndex
	return append(m[:insertionIndex], append([][]byte{row}, m[insertionIndex:]...)...), nil
}

// ReduceToIdentityMatrix takes a Matrix that is already in RREF form. It reduces the Matrix to the identity Matrix.
func ReduceRREFToIdentityMatrix(m Matrix) {
	var scale byte
	for i, row := range m {
		for j := i + 1; j < len(m); j++ {
			scale = m[i][j]
			galMulSliceXor(scale, m[j], row, &defaultOptions)
		}
	}
}

// MulVectorMatrix left-multiplies the matrix with the provided vector.
func MulVectorMatrix(v []byte, m Matrix) ([]byte, error) {
	if len(v) != len(m) {
		return nil, fmt.Errorf("vector length %d must match matrix height %d", len(v), len(m))
	}

	ret := make([]byte, len(m[0]))
	var scratch byte
	for i := 0; i < len(m[0]); i++ {
		scratch = 0
		for j, vectorVal := range v {
			scratch ^= galMultiply(m[j][i], vectorVal)
		}
		ret[i] = scratch
	}

	return ret, nil
}

func findLeadingEntryIndex(row []byte) (int, error) {
	for i, v := range row {
		if v != 0 {
			return i, nil
		}
	}

	return -1, fmt.Errorf("no leading entry found")
}
