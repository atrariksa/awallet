# AWALLET

This simple e-wallet uses some approach to optimize its apis performance such as :
1. implement master slave database (write and read only). this app has each connection for read and write.
2. implement redis as cache on layer repo read to decrease database load.
3. store and update balance on table user when TOPUP, INCOMING, OUTGOING happens. so that when api get balance called, it does not need to sum user mutation.
4. store and update user outgoing on its own database (user_total_outgoing) when transfer happens. so that no need to select all outgoing and aggregate them when api list top called. instead, when this api called, it just need to query select and order desc and limit which will have better performance than summing all outgoing and followed by aggregate by user.
5. add indexes

# How to run

1.  prepare database server and new database
2.  prepare redis server
3.  adjust config (.env) values for prepared database and redis
4.  run command with args "migrate up"
    example : go run main.go migrate up or ./awallet migrate up
5.  run command with args "server"
    example : go run main.go server or ./awallet server
6.  run command with args "blackbox" (run test cases)
    example : go run main.go blackbox or ./awallet blackbox

Before running command "blackbox", please make sure command "migrate up" and "server" already
done. so, apis can be tested by "blackbox". Please becareful, running command "blackbox" will clean up tables.