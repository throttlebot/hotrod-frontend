import requests, time, sys
from detect import detect

"""
Checks that the server returns a valid response
"""
def test(host, port):
	payload = {'customer': '392', 'nonse': '0.5'}
	response = requests.get("http://" + host + ":" + str(port) + "/dispatch", params=payload).json()
	assert response["ETA"] == 120000000000, "Incorrect Response {}".format(response)
	print response

if __name__ == '__main__':

	test(detect(sys.argv[1]), 8080)
