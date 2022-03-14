# Factorial RESTful API
Takes a POST Request that takes json with the following structure: {"a":#,"b":#} (where # is an int) and 
returns a json string containing the factorial of a and b in the same structure. When a negative number or 
other incorrect input is given, the API returns an error.