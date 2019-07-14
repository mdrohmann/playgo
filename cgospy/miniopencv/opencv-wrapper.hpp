#ifndef OPENCV_WRAPPER_HPP
#define OPENCV_WRAPPER_HPP

#ifdef __cplusplus
#include <opencv2/opencv.hpp>
#include <cstdlib>
extern "C"
{
  typedef cv::Mat* CvMatrix;
  typedef int32_t CInt;

#else
#include "stdlib.h"
typedef void* CvMatrix;
typedef int CInt;
#endif

  CvMatrix newCvMat();

  int cvMatAt(CvMatrix m, int x, int y);

  int captureImage(int device, CvMatrix edges);

  CInt* cvMatrixSize(CvMatrix m, int * len);

  void freeCvMat(CvMatrix mat);

  void imShow(CvMatrix mat);

#ifdef __cplusplus
}
#endif

#endif /* OPENCV_WRAPPER_HPP */
