pragma solidity 0.4.24;

contract B {
    function B() {}

    function bar() constant returns(uint) {
        return 100;
    }
}
