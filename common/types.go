package common

const (
	ErrorCodeSuccess            = 0
	ErrorCodeRequiredParamEmpty = 1
	ErrorCodeEmailExist         = 2
	ErrorCodeUnauthrized        = 3
	ErrorMerchantIdNotExist     = 4
	ErrorMerchantStillHasMember = 5
	ErrorEmailFormatInvalid     = 6
	ErrorRoleNotAvailable       = 7
	ErrorCodeUndefined          = 8
)

var ErrorMessageMap = map[int]string{
	ErrorCodeRequiredParamEmpty: "required parameter: %s can not be empty",
	ErrorCodeEmailExist:         "email has been used before",
	ErrorCodeUnauthrized:        "unauthorized",
	ErrorMerchantIdNotExist:     "merchant id does not exist",
	ErrorMerchantStillHasMember: "merchant still has some members",
	ErrorEmailFormatInvalid:     "invalid email format",
	ErrorRoleNotAvailable:       "Role is not available, only user, administrator, and superadmin allowed",
}

const (
	RoleUser          = "user"
	RoleAdministrator = "administrator"
	RoleSuperAdmin    = "superadmin"
)

var MapRoles = map[string]bool{
	RoleUser:          true,
	RoleAdministrator: true,
	RoleSuperAdmin:    true,
}

var DefaultMerchantPassword = "Merchant!234"

type ErrorMessage struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

type AuthenticationRequest struct {
	EmailAddress string
	Password     string
}

type MemberMerchant struct {
	MerchantCode string
}

type MemberRequest struct {
	Merchant     *MemberMerchant
	EmailAddress string
	Name         string
	Address      string
	Role         string
	Password     string
}

type MerchantRequest struct {
	Name    string
	Address string
}
