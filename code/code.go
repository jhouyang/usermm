package code

const (
    // CodeSucc          succ code
    CodeSucc            = 0

    // tcp 1000 ~ 2000
    // CodeTCPFailedGetUserInfo code succ
    CodeTCPFailedGetUserInfo    = 1101
    // CodeTCPPasswdErr password error
    CodeTCPPasswdErr            = 1102
    // CodeTCPInvalidToken invalid token
    CodeTCPInvalidToken         = 1200
    // CodeTCPTokenExpired token expired
    CodeTCPTokenExpired         = 1201
    // CodeTCPUserInfoNotMatch token info not match userinfo
    CodeTCPUserInfoNotMatch     = 1202
    // CodeTCPFailedUpdateUserInfo update userinfo failed
    CodeTCPFailedUpdateUserInfo = 1301
    // CodeTCPInternelErr internel error
    CodeTCPInternelErr          = 1401

    // HTTP 2000 ~ 3000
    // CodeInternalErr   internel err
    CodeInternalErr     = 2101
    // CodeTokenNotFound missing token
    CodeTokenNotFound   = 2102
    // CodeInvalidToken  token format is invalid
    CodeInvalidToken    = 2103
    // CodeErrBackend    failed to comm with backend server
    CodeErrBackend      = 2201
    // CodeInvalidPasswd passwd format isn't right
    CodeInvalidPasswd   = 2301
    // CodeFormFileFailed formFile get error
    CodeFormFileFailed  = 2401
    // CodeFileSizeErr file size not match (too small or too large)
    CodeFileSizeErr     = 2402
)

// CodeMsg code to msg description
var CodeMsg = map[int]string {
    // http
    CodeSucc          : "succ",
    CodeInternalErr   : "please try again!",
    CodeTokenNotFound : "param error: token not found",
    CodeInvalidToken  : "invalid token",
    CodeErrBackend    : "Error found!please try again!",
    CodeInvalidPasswd : "username/passwd error!",
    CodeFormFileFailed: "fetch file failed!",
    CodeFileSizeErr   : "File size err (should less than 5MB)!",

    // tcp
    CodeTCPFailedGetUserInfo    : "tcp server: failed to get userinfo",
    CodeTCPPasswdErr            : "tcp server: wrong passwd",
    CodeTCPInvalidToken         : "tcp server: invalid token format",
    CodeTCPTokenExpired         : "tcp server: token expired",
    CodeTCPUserInfoNotMatch     : "tcp server: token cache info not match",
    CodeTCPFailedUpdateUserInfo : "tcp server: failed to update userinfo",
    CodeTCPInternelErr          : "tcp server: internel error",
}
