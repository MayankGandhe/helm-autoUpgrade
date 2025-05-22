import os
import tempfile
import requests
from flask import Flask, request, jsonify
from threading import Thread
from helmpy import Helm

app = Flask(__name__)

# Helper to download a file from URL to a temp file
def download_to_temp(url, prefix):
    resp = requests.get(url)
    resp.raise_for_status()
    fd, path = tempfile.mkstemp(prefix=prefix)
    with os.fdopen(fd, 'wb') as f:
        f.write(resp.content)
    return path

def helm_upgrade_or_install(release_name, chart_url, values_url, namespace="default"):
    chart_path = download_to_temp(chart_url, "chart-")
    values_path = download_to_temp(values_url, "values-")
    helm = Helm()
    try:
        # Try upgrade first
        try:
            helm.upgrade(release_name, chart_path, namespace=namespace, values=values_path)
            print("Helm upgrade successful!")
        except Exception as e:
            if "not found" in str(e).lower():
                # If not found, do install
                helm.install(release_name, chart_path, namespace=namespace, values=values_path)
                print("Helm install successful!")
            else:
                print(f"Helm upgrade failed: {e}")
    finally:
        os.remove(chart_path)
        os.remove(values_path)

@app.route('/upgrade', methods=['POST'])
def upgrade_handler():
    data = request.get_json()
    release_name = data.get('releaseName')
    chart_url = data.get('chartURL')
    values_url = data.get('valuesURL')
    namespace = data.get('namespace', 'default')
    Thread(target=helm_upgrade_or_install, args=(release_name, chart_url, values_url, namespace)).start()
    return jsonify({"status": "acknowledged"}), 202

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
