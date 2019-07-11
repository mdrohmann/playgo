#include "opencv-wrapper.hpp"

extern "C"
{

  CvMatrix newCvMat() { return new cv::Mat(); }

  int cvMatAt(CvMatrix m, int x, int y) {
    return m->at<uint8_t>(cv::Point2i(x, y));
  }

  int* cvMatrixSize(CvMatrix m, int* len)
  {
    cv::MatSize mSize = m->size;
    *len = m->dims;
    return mSize.p;
  }

  int captureImage(int device, CvMatrix edges)
  {
    cv::VideoCapture cap(device);
    if (!cap.isOpened()) {
      return -1;
    }

    cv::Mat frame;
    cap >> frame; // get a new frame from camera
    cv::cvtColor(frame, *edges, cv::COLOR_BGR2GRAY);
    cv::GaussianBlur(*edges, *edges, cv::Size(7, 7), 1.5, 1.5);
    cv::Canny(*edges, *edges, 0, 30, 3);

    return 0;
  }

  void freeCvMat(CvMatrix mat)
  {
    if (mat != nullptr) {
      delete mat;
    }
  }

  void imShow(CvMatrix mat)
  {
    cv::imshow("cheese", *mat);
    cv::waitKey(0);
    cv::destroyAllWindows();
  }
}