npx prettier -w src/*.js
docker build -t docker.io/cmwylie19/find-me .; docker push docker.io/cmwylie19/find-me; k rollout restart deploy/find-me
