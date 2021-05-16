"""Python testing harness"""

import logging
logging.basicConfig(level=logging.DEBUG)

import requests

LOGGER = logging.getLogger(__name__)

def upload_attachment(job_id: str, file_path: str):
    """Function used to upload attachment
    for a given job"""

    with open(file_path, 'r') as f:
        data = f.read()
    url = 'http://localhost:10312/jobs/{}/attachments'.format(job_id)
    try:
        r = requests.post(url, headers={'X-Authenticated-Userid': 'psauerborn'},
                          files={'attachment': data})
        LOGGER.debug('received API response %s', r.text)
        r.raise_for_status()

    except requests.HTTPError:
        LOGGER.exception('unable to upload attachment')
        raise

if __name__ == '__main__':

    job_id = 'c4ff0853-836f-4354-91d3-6614df17d02b'
    upload_attachment(job_id, 'tests/test.txt')