package repository

type ReportRepo interface {
	GetReportByID()
	GetReport()
	UpdateReport() error
	CreateReport() error
	DeleteReport() error
}

func (r repo) GetReportByID() {

}

func (r repo) GetReport() {

}

func (r repo) UpdateReport() error {
	return nil
}

func (r repo) CreateReport() error {
	return nil
}

func (r repo) DeleteReport() error {
	return nil
}
