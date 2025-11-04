import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { AuthService } from "./api/v1/auth_pb";
import { APIKeyService } from "./api/v1/apikey_pb";

/**
 * API Client for making authenticated requests to the backend
 * Automatically includes JWT token in Authorization header
 */
class ApiClient {
  private transport;
  private authClient: ReturnType<typeof createClient<typeof AuthService>>;
  private apiKeyClient: ReturnType<typeof createClient<typeof APIKeyService>>;

  constructor() {
    // Create transport for Connect RPC
    // API is served on the same site, so we use window.location.origin
    this.transport = createConnectTransport({
      baseUrl: window.location.origin,
      // Interceptor to add JWT token to all requests
      interceptors: [
        (next) => async (req) => {
          const token = this.getToken();
          if (token) {
            req.header.set("Authorization", `Bearer ${token}`);
          }
          return await next(req);
        },
      ],
    });

    // Create auth service client
    this.authClient = createClient(AuthService, this.transport);

    // Create API key service client
    this.apiKeyClient = createClient(APIKeyService, this.transport);
  }

  /**
   * Get stored JWT token from localStorage
   */
  private getToken(): string | null {
    return localStorage.getItem("auth_token");
  }

  /**
   * Store JWT token in localStorage
   */
  setToken(token: string): void {
    localStorage.setItem("auth_token", token);
  }

  /**
   * Remove JWT token from localStorage
   */
  clearToken(): void {
    localStorage.removeItem("auth_token");
  }

  /**
   * Get the auth service client
   */
  get auth() {
    return this.authClient;
  }

  /**
   * Get the API key service client
   */
  get apiKey() {
    return this.apiKeyClient;
  }
}

// Export singleton instance
export const apiClient = new ApiClient();
