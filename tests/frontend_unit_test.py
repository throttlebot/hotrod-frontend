import requests, time, sys
from detect import detect

command = "kubectl get service hotrod-frontend -o jsonpath='{.status.loadBalancer.ingress[0].ip}' -n "

"""
Checks that the server returns a valid response
"""
def test(host, port):
	response = requests.get("http://" + host + ":" + str(port))
	assert response.status_code == 200, "Status code was {} instead of 200".format(response.status_code)

if __name__ == '__main__':

	test("hotrod-frontend", 8080)
