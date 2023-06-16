Pre-Reqs: Docker
Step 1: Build the image

cd vlg/dockerize_demo
docker build --no-cache -t image_name .


Step 2: Run the container. Make sure to forward 8888 since the notebook will be running from the container
docker run -p 8888:8888 image_name
