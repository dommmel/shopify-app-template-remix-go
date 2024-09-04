
export const loader = async () => {
  return renderAppBridge( );
}

function renderAppBridge() {
  const html = `
  <head>
      <meta name="shopify-api-key" content="${process.env.SHOPIFY_API_KEY}" />
      <script src="https://cdn.shopify.com/shopifycloud/app-bridge.js"></script>
  </head>`

  throw new Response(html, {
    headers: {
      'Content-Type': 'text/html;charset=utf-8',
    },
  })
}

