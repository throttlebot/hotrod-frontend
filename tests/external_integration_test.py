import requests, time, sys

"""
Checks that the server returns a valid response
"""
def test(host, port):
	payload = {'customer': '392', 'nonse': '0.5'}
	response = requests.get("http://" + host + ":" + str(port) + "/dispatch", params=payload).json()
	assert isinstance(response["ETA"], int), "Incorrect Response {}".format(response)
	print response

if __name__ == '__main__':

	test(sys.argv[1], 80)
