These are instructions on how to generate the needed credentials for this receptor. 

GitLab receptor needs a GroupID and a Token to access an account.


First, get the group ID
1. In the top-left corner, select **Menu**.
2. Select **Groups**.
3. Select **Your groups**.
4. Select a group.
5. Copy the `group ID` under the **group name**.
6. Paste the `group ID` into the **Group ID** field.

Next, create a read-only personal access token
1. In the top-right corner, select your avatar.
2. Select **Edit** profile.
3. On the left sidebar, select **Access Tokens**.
4. Enter a name and optional expiry date for the token.
5. Select `read_api`, `read_user`, and `read_repository` scopes.
6. Select **Create personal access token**.
7. Enter the personal access token into the **Token** field.