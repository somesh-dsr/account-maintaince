# account-maintaince
User has to run initiate-server.go file to start API server.
By default server will listen on port 100 with localhost.
By default API server will create user database in the tmp folder

Note:Expiry we have to give in a unix time stamp format. To get Unix time for a particuar time period follow below link for reference.

https://www.epochconverter.com/

API details:


127.0.0.1:100/users/account/credit  - to credit amount from a particular account
127.0.0.1:100/users/account/debit   - to debit amount from a particular account
127.0.0.1:100/users/create          - to create a new user acccount with all the details information
127.0.0.1:100/users/expired/amount?userId=<user id/ uuid>  - to get expired amount details
127.0.0.1:100/users/debit/logs?userId=<user id/ uuid>      - to get debited amount logs for a particular user
127.0.0.1:100/users/credit/logs?userId=<user id/ uuid>     - to get credited amount logs for a particular user


Sample Json for above APIS:

Account Creation:

{
"name":"somesh",
"mailId" : "dsr@gmail.com",
"mobileNo" : "9743761740",
"address":{
"pincode" : 500032,
"street" :   "Gachibowli",
"city" : "Hyderabad",
"state": "Telangana"
}
}

Output Json with account UUID:

{
    "statusCode": 201,
    "success": "Account Created Successfully: key is 830e157cbf13f82fbc24"
}

Account credit with specified amount:

{
 "activity":"credit",
 "payload":{
 "userId":"3049ae8d4c41b6dd1a67",
 "amount":100,
 "type":"subscription",
 "priority":10,
 "expiry":16
 }
}

output Json:

{
    "statusCode": 201,
    "success": "Amount Credited successfully"
}

Account Debit with specified amount:
{
"activity":"debit",
"payload": {
 "userId":"3049ae8d4c41b6dd1a67",
 "amount":21
}
}

output Json:
{
    "statusCode": 201,
    "success": "Amount Debited successfully"
}


Sample Json Output for credit,debit and expired amount logs:

[
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 20:46:54.059994272 +0530 IST m=+39.700318378"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 20:51:39.525332815 +0530 IST m=+5.989051424"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 20:52:43.72874011 +0530 IST m=+70.192458720"
    },
    {
        "amount": 21,
        "amountPriority": 0,
        "transitionDate": "2021-03-10 20:54:39.698486551 +0530 IST m=+186.162205236"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 21:27:49.644338082 +0530 IST m=+47.737393516"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 21:30:16.3688149 +0530 IST m=+7.053497363"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 21:31:54.926208123 +0530 IST m=+4.915768743"
    },
    {
        "amount": 21,
        "amountPriority": 0,
        "transitionDate": "2021-03-10 21:34:15.691846098 +0530 IST m=+145.681406883"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 21:51:54.378332071 +0530 IST m=+37.285154854"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 21:54:57.104017485 +0530 IST m=+4.047966078"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 21:57:11.37146589 +0530 IST m=+11.191714669"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 21:59:09.925027297 +0530 IST m=+5.728815502"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 22:00:42.597504624 +0530 IST m=+3.148552758"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 22:03:14.247564816 +0530 IST m=+2.998607203"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 22:14:35.098927856 +0530 IST m=+683.849970253"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-10 22:15:18.221053111 +0530 IST m=+726.972095495"
    },
    {
        "amount": 21,
        "amountPriority": 0,
        "transitionDate": "2021-03-10 22:16:42.1600104 +0530 IST m=+810.911052945"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-11 19:58:49.96796774 +0530 IST m=+28.607301717"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-11 20:03:33.555985203 +0530 IST m=+27.230963169"
    },
    {
        "amount": 100,
        "amountPriority": 10,
        "transitionDate": "2021-03-11 20:04:44.12562788 +0530 IST m=+17.528908339"
    }
]


Some failure Json output:

User ceation failed:

{
    "error": "UNIQUE constraint failed: usersaccount.mobile",
    "statusCode": 409
}

Account Debit/credit:

{
    "error": "Account not found.",
    "statusCode": 409
}
