USERPOOLID=eu-west-xxxx; \
for u in $(aws --profile prod cognito-idp list-users --user-pool $USERPOOLID --filter 'cognito:user_status = "RESET_REQUIRED"' | jq -r '.Users[].Username'); \
do aws --profile prod cognito-idp admin-delete-user --user-pool-id $USERPOOLID --username $u; \
done