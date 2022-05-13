# running instructions:
# from root of the project, execute me as
# `. ./scripts/build_image.sh`
# `. ./scripts/build_image.sh v1.6.0`

# set required environment variables
REMINDER_IMAGE=reminder
REMINDER_CONTAINER=reminder
TAG="${@:-latest}"

# building and tag image
echo "STEP-01: Build the latest image"
docker build -t goyalmunish/${REMINDER_IMAGE} -f Dockerfile ./
docker tag goyalmunish/${REMINDER_IMAGE} goyalmunish/${REMINDER_IMAGE}:${TAG}
echo "STEP-03: Push the images"
docker push goyalmunish/${REMINDER_IMAGE}:latest
docker push goyalmunish/${REMINDER_IMAGE}:${TAG}
echo "STEP-03: Pull the default image"
docker pull goyalmunish/${REMINDER_IMAGE}:latest
docker pull goyalmunish/${REMINDER_IMAGE}:${TAG}
