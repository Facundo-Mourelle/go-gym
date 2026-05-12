import type { AxiosError, AxiosInstance, InternalAxiosRequestConfig } from 'axios';
import axios from 'axios';

const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

class ApiClient {
    private client: AxiosInstance;

    constructor() {
        this.client = axios.create({
            baseURL: BASE_URL,
            headers: {
                'Content-Type': 'application/json',
            },
        });

        // Request interceptor - add auth token
        this.client.interceptors.request.use(
            (config: InternalAxiosRequestConfig) => {
                const token = localStorage.getItem('auth_token');

                // DETAILED DEBUG LOGGING
                console.group('🔐 API Request');
                console.log('URL:', config.url);
                console.log('Method:', config.method);
                console.log('Token in localStorage:', token ? `${token.substring(0, 20)}...` : 'NOT FOUND');
                console.log('Headers before:', config.headers);

                if (token && config.headers) {
                    config.headers.Authorization = `Bearer ${token}`;
                    console.log('✅ Authorization header added');
                } else {
                    console.log('❌ No token or no headers object');
                }

                console.log('Headers after:', config.headers);
                console.groupEnd();

                return config;
            },
            (error) => {
                return Promise.reject(error);
            }
        );

        // Response interceptor - handle errors
        this.client.interceptors.response.use(
            (response) => response,
            (error: AxiosError) => {
                console.error('❌ API Error:', error.response?.status, error.response?.data);

                if (error.response?.status === 401) {
                    console.log('🔓 Unauthorized - clearing auth and redirecting');
                    localStorage.removeItem('auth_token');
                    localStorage.removeItem('auth-storage');
                    window.location.href = '/login';
                }
                return Promise.reject(error);
            }
        );
    }

    getInstance(): AxiosInstance {
        return this.client;
    }
}

export const apiClient = new ApiClient().getInstance();
