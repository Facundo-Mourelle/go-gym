export interface User {
    user_id: string;
    email: string;
    name: string;
}

export interface AuthResponse {
    token: string;
    user_id: string;
    email: string;
    name: string;
}

export interface LoginRequest {
    email: string;
    password: string;
}

export interface RegisterRequest {
    email: string;
    password: string;
    name: string;
}
