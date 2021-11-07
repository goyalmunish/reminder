# running instructions:
# from root of the project, execute me as
# `. ./scripts/build_image.sh`

# set required environment variables
REMINDER_IMAGE=reminder
REMINDER_CONTAINER=reminder
TAG="${@:-latest}"

# building and pushing default (bionic) image
echo "STEP-01: Build and tagging the default image"
docker build -t goyalmunish/${REMINDER_IMAGE} -f Dockerfile ./
echo "STEP-02: Push the default image"
docker push goyalmunish/${REMINDER_IMAGE}:${TAG}
echo "STEP-03: Pull the default image"
docker pull goyalmunish/${REMINDER_IMAGE}:${TAG}
