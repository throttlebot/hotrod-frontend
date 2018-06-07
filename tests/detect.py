def detect(namespace):
    while True:
        value = subprocess.check_output(command + namespace, shell=True).strip()
        assert value != "<none>", "hortod must have an external ip"
        if value == "<pending>":
            time.sleep(1)
        else:
            return value
