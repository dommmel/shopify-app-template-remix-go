import { TitleBar } from "@shopify/app-bridge-react";
import {
  Layout,
  Page,
  Card,
} from "@shopify/polaris";

export default function Index() {
  return (
    <Page>
      <TitleBar title="Additional page" />
      <Layout>
        <Layout.Section>
        <Card>
          Nothing to see yet
        </Card>
        </Layout.Section>
      </Layout>
    </Page>
    
  )
}