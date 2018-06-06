import requests, subprocess, time, sys

command = "kubectl get service hotrod-frontend -o jsonpath='{.status.loadBalancer.ingress[0].ip}' -n "

"""
Checks that the server returns a valid response
"""
def test(host, port):
	response = requests.get("http://" + host + ":" + str(port))
	assert response.status_code == 200, "Status code was {} instead of 200".format(response.status_code)

        payload = {'customer': '392', 'nonse': '0.5'}
	response = requests.get("http://" + host + ":" + str(port) + "/dispatch", params=payload).json()
	assert response["ETA"] == 120000000000, "Incorrect Response {}".format(response)
	print response

def detect(namespace):
	while True:
		value = subprocess.check_output(command + namespace, shell=True).strip()
		assert value != "<none>", "hortod must have an external ip"
		if value == "<pending>":
			time.sleep(1)
		else:
			return value

if __name__ == '__main__':

	test(detect(sys.argv[1]), 8080)
