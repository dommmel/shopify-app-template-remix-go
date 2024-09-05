## Ideas
- Merchant impersonation
  - Use the app on a whitelisted admin store (e.g. my-app-admin.myshopfy.com)
  - Check a `?impersonate=store-to-impersonate.myshopify.com` query parameter and in `FindOrCreateUserByEncodedSessionToken`check if the logged-in store is the whitelisted admin store and if so return the data for the store given in the query string.
  - This should work as long as everything Shopfiy-related is done in the backend so we can swap out the content that the backend is deliviering while being logged in as another user (admin store).
- Admin route with separate Auth, for viewing user records and linking to the above
