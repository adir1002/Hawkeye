from elasticsearch import Elasticsearch

es = Elasticsearch('http://elastic-hawkeye.westeurope.cloudapp.azure.com:9200/')

from datetime import datetime
import time
import random

id_input = 0
normal_rounds_counter = 0
while(True):
    if (normal_rounds_counter < 3):
        es = Elasticsearch('http://elastic-hawkeye.westeurope.cloudapp.azure.com:9200/')
        logs_counter = 0
        logs_amount = random.randint(100, 200)
        while(logs_counter < logs_amount):
            es.index(index="logs", doc_type="test-type", id=id_input, body={"any": "data", "timestamp": datetime.utcnow()})
            {u'_id': id_input, u'_index': u'my-index-logs', u'_type': u'test-type', u'_version': 1, u'ok': True}
            id_input += 1
            logs_counter += 1
        normal_rounds_counter += 1
        time.sleep(random.randint(240, 300))
    else:
        es = Elasticsearch('http://elastic-hawkeye.westeurope.cloudapp.azure.com:9200/')
        logs_counter = 0
        logs_amount = random.randint(500, 1000)
        while(logs_counter < logs_amount):
            es.index(index="logs", doc_type="test-type", id=id_input, body={"any": "data", "timestamp": datetime.utcnow()})
            {u'_id': id_input, u'_index': u'my-index-logs', u'_type': u'test-type', u'_version': 1, u'ok': True}
            id_input += 1
            logs_counter += 1
        normal_rounds_counter = 0
        time.sleep(random.randint(240, 300))
