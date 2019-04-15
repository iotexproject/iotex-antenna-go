pragma solidity 0.4.24;

contract B {
    address aa;
    function B(address a) {
        aa=a;
    }
    function bar(uint x) constant returns(uint) {
        return x;
    }
    function barstring(string y)returns(string){
        return y;
    }
    function getaddress()returns(address){
            return aa;
    }
}
