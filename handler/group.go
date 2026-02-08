package handler

type HandlerGroup struct {
	Product     *ProductHandler
	Category    *CategoryHandler
	Transaction *TransactionHandler
}
