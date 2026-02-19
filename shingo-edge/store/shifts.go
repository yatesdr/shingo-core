package store

// Shift represents a work shift with start/end times.
type Shift struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ShiftNumber int    `json:"shift_number"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

// ListShifts returns all shifts ordered by shift_number.
func (db *DB) ListShifts() ([]Shift, error) {
	rows, err := db.Query("SELECT id, name, shift_number, start_time, end_time FROM shifts ORDER BY shift_number")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shifts []Shift
	for rows.Next() {
		var s Shift
		if err := rows.Scan(&s.ID, &s.Name, &s.ShiftNumber, &s.StartTime, &s.EndTime); err != nil {
			return nil, err
		}
		shifts = append(shifts, s)
	}
	return shifts, rows.Err()
}

// UpsertShift inserts or replaces a shift by shift_number.
func (db *DB) UpsertShift(shiftNumber int, name, startTime, endTime string) error {
	_, err := db.Exec(
		`INSERT INTO shifts (shift_number, name, start_time, end_time)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(shift_number) DO UPDATE SET name=excluded.name, start_time=excluded.start_time, end_time=excluded.end_time`,
		shiftNumber, name, startTime, endTime,
	)
	return err
}

// DeleteShift removes a shift by shift_number.
func (db *DB) DeleteShift(shiftNumber int) error {
	_, err := db.Exec("DELETE FROM shifts WHERE shift_number = ?", shiftNumber)
	return err
}
