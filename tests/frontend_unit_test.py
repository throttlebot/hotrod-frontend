import requests, subprocess, time, sys

command = "kubectl get service hotrod-frontend -o jsonpath='{.status.loadBalancer.ingress[0].ip}' -n "

"""
Checks that the server returns a valid response
"""
def test(host, port):
	response = requests.get("http://" + host + ":" + str(port))
	assert response.status_code == 200, "Status code was {} instead of 200".format(response.status_code)

if __name__ == '__main__':

	test(detect(sys.argv[1]), 8080)
