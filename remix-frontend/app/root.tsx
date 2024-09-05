import { json } from "@remix-run/node";
import { Link,Links,Meta,  Scripts, ScrollRestoration, Outlet, useLoaderData } from "@remix-run/react";
import {AppProvider} from '@shopify/shopify-app-remix/react';
import { NavMenu } from "@shopify/app-bridge-react";
import polarisStyles from "@shopify/polaris/build/esm/styles.css?url";

export const links = () => [{ rel: "stylesheet", href: polarisStyles }];

export const loader = async () => {
  return json({ apiKey:  "fa6737ec2ded9ca5ddfbd17fc0e3450b" });
};

export default function App() {
  const { apiKey } = useLoaderData<typeof loader>();

  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body>
      <AppProvider isEmbeddedApp apiKey={apiKey}>
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

