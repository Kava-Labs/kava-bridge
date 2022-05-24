package types

type TssMsg interface {
	ValidateBasic() error
}
