#ifndef OPENCV_WRAPPER_HPP
#define OPENCV_WRAPPER_HPP

#ifdef __cplusplus
#include <opencv2/opencv.hpp>
extern "C"
{
  typedef cv::Mat* CvMatrix;

#else
typedef void* CvMatrix;
#endif

  CvMatrix newCvMat();

  int captureImage(int device, CvMatrix edges);

  void freeCvMat(CvMatrix mat);

  void imShow(CvMatrix mat);

#ifdef __cplusplus
}
#endif

#endif /* OPENCV_WRAPPER_HPP */
