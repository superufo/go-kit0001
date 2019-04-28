package service


service HelloService {
	rpc Hello (String) returns (String);

	rpc Channel (stream String) returns (stream String)
}

