package model

import (
	"fmt"
	"regexp"
)

// ValidateInput valida un valore in ingresso secondo i vincoli specificati.
// Parametri:
//   - value: il valore da validare (interface{} per poter verificare il tipo)
//   - min: numero minimo di caratteri
//   - max: numero massimo di caratteri
//   - pattern: espressione regolare che il valore deve rispettare
//   - tipo: tipo atteso (es. "string")
//
// Ritorna un errore descrittivo se la validazione fallisce, nil se il valore è valido.
func ValidateInput(value interface{}, min int, max int, pattern string, tipo string) error {

	// Default negativo: l'input è considerato invalido fino a prova contraria
	var err error

	// 1. Verifica del tipo
	switch tipo {
	case "string":
		str, ok := value.(string)
		if !ok {
			err = fmt.Errorf("expected type 'string', got a different type")
		} else if len(str) < min {
			// 2. Verifica lunghezza minima
			err = fmt.Errorf("input is too short: minimum %d characters required, got %d", min, len(str))
		} else if len(str) > max {
			// 3. Verifica lunghezza massima
			err = fmt.Errorf("input is too long: maximum %d characters allowed, got %d", max, len(str))
		} else {
			// 4. Verifica regex
			re, compileErr := regexp.Compile(pattern)
			if compileErr != nil {
				err = fmt.Errorf("invalid regex pattern '%s': %w", pattern, compileErr)
			} else if !re.MatchString(str) {
				err = fmt.Errorf("input '%s' does not match the required pattern '%s'", str, pattern)
			} else {
				// Tutto ok: l'input rispetta tutti i vincoli
				err = nil
			}
		}
	default:
		err = fmt.Errorf("unsupported type '%s'", tipo)
	}

	return err
}
