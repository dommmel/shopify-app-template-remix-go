import { Link,Links,Meta,  Scripts, ScrollRestoration, Outlet } from "@remix-run/react";
import {AppProvider} from '@shopify/shopify-app-remix/react';
import { NavMenu } from "@shopify/app-bridge-react";
import polarisStyles from "@shopify/polaris/build/esm/styles.css?url";

export const links = () => [{ rel: "stylesheet", href: polarisStyles }];

export default function App() {
  const { SHOPIFY_API_KEY } = import.meta.env;
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="viewport" content="width=device-width,initial-scale=1" />
        <link rel="preconnect" href="https://cdn.shopify.com/" />
        <link
          rel="stylesheet"
          href="https://cdn.shopify.com/static/fonts/inter/v4/styles.css"
        />
        <Meta />
        <Links />
      </head>
      <body>
      <AppProvider isEmbeddedApp apiKey={SHOPIFY_API_KEY}>
        <NavMenu>
          <Link to="/" rel="home">
            Home
          </Link>
          <Link to="/additional">Additional page</Link>
        </NavMenu>
        <Outlet />
        <ScrollRestoration />
        <Scripts />
      </AppProvider>
      </body>
    </html>
  );
}