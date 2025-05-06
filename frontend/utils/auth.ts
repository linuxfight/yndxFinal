export const useAuth = () => {
    const cookieName = 'app_auth_token'
    const cookieOptions = {
        path: '/',
        secure: true,
        sameSite: 'strict' as const,
        maxAge: 86400, // 24 hour expiration
        httpOnly: false // Important: Only set to true if set via server-side
    }

    // Check if token exists
    const checkTokenExists = (): boolean => {
        return !!useCookie(cookieName).value
    }

    // Get token value
    const getToken = (): string | null => {
        return useCookie(cookieName).value || null
    }

    // Set token with security options
    const setToken = (token: string): void => {
        const authCookie = useCookie(cookieName, cookieOptions)
        authCookie.value = token
    }

    // Delete token
    const deleteToken = (): void => {
        const authCookie = useCookie(cookieName, {
            ...cookieOptions,
            maxAge: -1 // Expire immediately
        })
        authCookie.value = null
    }

    return {
        checkTokenExists,
        getToken,
        setToken,
        deleteToken
    }
}