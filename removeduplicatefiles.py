import os
import hashlib

def md5sum(filename):
    hash_md5 = hashlib.md5()
    with open(filename, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
            hash_md5.update(chunk)
    return hash_md5.hexdigest()

seen_hashes = set()
current_dir = '.'

for filename in os.listdir(current_dir):
    if os.path.isfile(filename):
        file_path = os.path.join(current_dir, filename)
        file_hash = md5sum(file_path)
        
        if file_hash in seen_hashes:
            os.remove(file_path)
        else:
            seen_hashes.add(file_hash)

