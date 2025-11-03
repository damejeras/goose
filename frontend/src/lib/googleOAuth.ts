/**
 * Google OAuth Configuration
 *
 * To set up Google OAuth:
 * 1. Go to https://console.cloud.google.com/
 * 2. Create a new project or select an existing one
 * 3. Enable Google+ API
 * 4. Go to Credentials > Create Credentials > OAuth 2.0 Client ID
 * 5. Configure authorized JavaScript origins (e.g., http://localhost:3000)
 * 6. Replace GOOGLE_CLIENT_ID below with your actual client ID
 */

// TODO: Replace with your actual Google Client ID
export const GOOGLE_CLIENT_ID =
  "815902425628-ndl0rlr1k5c6r7iqk85147tacbujp5do.apps.googleusercontent.com";

/**
 * Initialize Google Sign-In
 * This should be called once when the app loads
 */
export function initGoogleSignIn(): Promise<void> {
  return new Promise((resolve, reject) => {
    // Check if gsi script is already loaded
    if (window.google?.accounts?.id) {
      resolve();
      return;
    }

    // Load Google Sign-In script
    const script = document.createElement("script");
    script.src = "https://accounts.google.com/gsi/client";
    script.async = true;
    script.defer = true;
    script.onload = () => resolve();
    script.onerror = () =>
      reject(new Error("Failed to load Google Sign-In script"));
    document.head.appendChild(script);
  });
}

/**
 * Trigger Google OAuth popup and return the ID token
 */
export function signInWithGoogle(): Promise<string> {
  return new Promise((resolve, reject) => {
    if (!window.google?.accounts?.id) {
      reject(new Error("Google Sign-In not initialized"));
      return;
    }

    // Initialize Google Sign-In with callback
    window.google.accounts.id.initialize({
      client_id: GOOGLE_CLIENT_ID,
      callback: (response: { credential: string }) => {
        // response.credential contains the JWT ID token
        resolve(response.credential);
      },
    });

    // Trigger the One Tap prompt
    window.google.accounts.id.prompt((notification) => {
      if (notification.isNotDisplayed() || notification.isSkippedMoment()) {
        // One Tap was not displayed or skipped, fall back to button click
        // This is handled by the login button in the UI
        console.log(
          "One Tap not available:",
          notification.getNotDisplayedReason(),
        );
      }
    });
  });
}

/**
 * Initialize Google Sign-In with callback for custom button
 * @param onSuccess - Callback when sign-in succeeds
 */
export function initializeGoogleCallback(
  onSuccess: (idToken: string) => void,
): void {
  if (!window.google?.accounts?.id) {
    throw new Error("Google Sign-In not initialized");
  }

  // Initialize with callback
  window.google.accounts.id.initialize({
    client_id: GOOGLE_CLIENT_ID,
    callback: (response: { credential: string }) => {
      onSuccess(response.credential);
    },
  });
}

// Store the hidden Google button element
let hiddenGoogleButton: HTMLElement | null = null;

/**
 * Initialize Google Sign-In with a hidden button for popup flow
 * This renders an invisible Google button that we can trigger programmatically
 * @param onSuccess - Callback when sign-in succeeds
 * @param onError - Callback when sign-in fails
 */
export function initializeOAuth2Client(
  onSuccess: (idToken: string) => void,
  onError?: (error: Error) => void,
): void {
  if (!window.google?.accounts?.id) {
    throw new Error("Google Sign-In not initialized");
  }

  // Initialize with callback to get ID token
  window.google.accounts.id.initialize({
    client_id: GOOGLE_CLIENT_ID,
    callback: (response: { credential: string }) => {
      // response.credential is the ID token (JWT format)
      onSuccess(response.credential);
    },
  });

  // Create a hidden div for the Google button
  hiddenGoogleButton = document.createElement("div");
  hiddenGoogleButton.style.position = "fixed";
  hiddenGoogleButton.style.top = "-9999px";
  hiddenGoogleButton.style.left = "-9999px";
  hiddenGoogleButton.style.visibility = "hidden";
  document.body.appendChild(hiddenGoogleButton);

  // Render the Google button (hidden)
  try {
    window.google.accounts.id.renderButton(hiddenGoogleButton, {
      theme: "outline",
      size: "large",
      text: "signin_with",
    });
  } catch (error) {
    if (onError) {
      onError(
        error instanceof Error ? error : new Error("Failed to render button"),
      );
    }
  }
}

/**
 * Trigger Google Sign-In popup by clicking the hidden Google button
 */
export function triggerGoogleSignIn(): void {
  if (!hiddenGoogleButton) {
    throw new Error("Google button not initialized");
  }

  // Find and click the actual Google button inside the hidden div
  const googleBtn = hiddenGoogleButton.querySelector("div[role='button']");
  if (googleBtn) {
    (googleBtn as HTMLElement).click();
  } else {
    throw new Error("Google button not found");
  }
}

// Type definitions for Google Sign-In
declare global {
  interface Window {
    google?: {
      accounts: {
        id: {
          initialize: (config: {
            client_id: string;
            callback: (response: { credential: string }) => void;
          }) => void;
          prompt: (
            callback?: (notification: {
              isNotDisplayed: () => boolean;
              isSkippedMoment: () => boolean;
              getNotDisplayedReason: () => string;
            }) => void,
          ) => void;
          renderButton: (
            element: HTMLElement,
            options: {
              theme?: "outline" | "filled_blue" | "filled_black";
              size?: "large" | "medium" | "small";
              text?: "signin_with" | "signup_with" | "continue_with" | "signin";
              width?: number;
            },
          ) => void;
        };
      };
    };
  }
}
