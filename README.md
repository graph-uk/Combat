# Combat
Test framework, where each test is independent application

Getting Combat on Windows
=====
* Install git
* Install go. For example, you've set gopath as `C:\gopath`
* Run `go get github.com/graph-uk/Combat`
* Add go packages binary `C:\gopath\bin` to path.
* Add your test's directory to GOPATH. If you are has not Combat tests, you are able add directory with default example tests.
`C:\gopath\src\github.com\graph-uk\Combat\Tests_Examples`

Check your installation
=====
Go to test directory, run `Combat list`. Now you should see information about test that you have.
```
simpleTestReturnsTrue
-------------------------------------------------
Locale               EnumParam           EN RU
AdminName            StringParam
HostName             StringParam
SessionTimestamp     StringParam
simpleTestReturnsFalse
-------------------------------------------------
Resolution           EnumParam           DesktopView MobileView
HostName             StringParam
SessionTimestamp     StringParam
Locale               EnumParam           EN RU US
```
Now you are able add your own test directory to GOPATH, and write your own tests.

Get addition information
=====
Run `Combat help`