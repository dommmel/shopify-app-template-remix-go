// copy and pasted from https://github.com/Shopify/shopify-app-js/blob/4ed3fe6753c51e155ea687c0336b0804d656a260/packages/apps/shopify-app-remix/src/server/authenticate/admin/authenticate.ts
import {UserResponse } from "../grpc/generated/user";
import { findOrCreateUserByEncodedSessionToken } from "../grpc/user.server";
import {redirect} from '@remix-run/server-runtime';


export function authenticate(request: Request): Promise<UserResponse> {
  const headerSessionToken = getSessionTokenHeader(request);
  const searchParamSessionToken = getSessionTokenFromUrlParam(request)
  const sessionToken = (headerSessionToken || searchParamSessionToken)!
  if (!sessionToken) {
    throw redirectToSessionTokenBouncePage(request);
  }

  return findOrCreateUserByEncodedSessionToken(sessionToken);

}

// Content of https://github.com/Shopify/shopify-app-js/blob/main/packages/apps/shopify-app-remix/src/server/authenticate/helpers/get-session-token-header.ts
const SESSION_TOKEN_PARAM = 'id_token';

export function getSessionTokenHeader(request: Request): string | undefined {
  return request.headers.get('authorization')?.replace('Bearer ', '');
}

export function getSessionTokenFromUrlParam(request: Request): string | null {
  const url = new URL(request.url);

  return url.searchParams.get(SESSION_TOKEN_PARAM);
}


// Inspired by https://github.com/Shopify/shopify-app-js/blob/main/packages/apps/shopify-app-remix/src/server/authenticate/admin/helpers/redirect-to-bounce-page.ts
// and https://github.com/Shopify/example-app--embedded-auth--js/blob/main/routes/auth.js
function redirectToSessionTokenBouncePage(request: Request) {
  const url = new URL(request.url);
  const searchParams = url.searchParams
  // Remove `id_token` from the query string to prevent an invalid session token sent to the redirect path.
  searchParams.delete('id_token')

  // Using shopify-reload path to redirect the bounce automatically.
  searchParams.append('shopify-reload',`${url.pathname}?${searchParams.toString()}`)
  return redirect(`/session-token-bounce?${searchParams.toString()}`)
}
