// Copy and pasted from: https://github.com/Shopify/shopify-app-js/blob/4ed3fe6753c51e155ea687c0336b0804d656a260/packages/apps/shopify-app-remix/src/server/authenticate/helpers/add-response-headers.ts

export function addDocumentResponseHeaders(
  request: Request,
  headers: Headers,
) {
  const {searchParams} = new URL(request.url);
  const shop = searchParams.get('shop')!
  if (shop) {
    headers.set(
      'Link',
      '<https://cdn.shopify.com/shopifycloud/app-bridge.js>; rel="preload"; as="script";',
    );
    headers.set(
      'Content-Security-Policy',
      `frame-ancestors https://${shop} https://admin.shopify.com https://*.spin.dev;`,
    );
  }
}