package core

type Transaction struct {
	Sender    string
	Receiver  string
	Amount    float64
	Timestamp int64
}

func NewTransaction(sender, receiver string, amount int) *Transaction {
	return &Transaction{Sender: sender, Receiver: receiver, Amount: float64(amount)}
}
