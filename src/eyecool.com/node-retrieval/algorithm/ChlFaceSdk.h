////////////////////////////////////////////////////////
//
// ChlFaceSdk.h
//
// 多通道人脸识别SDK接口定义
//
////////////////////////////////////////////////////////

#ifndef __CHLFACESDK__
#define __CHLFACESDK__

#include "typedef.h"

#ifdef _WIN32
  #ifdef IDFACESDK_EXPORTS
    #define CHLFACESDK_API __declspec(dllexport)
  #else
    #define CHLFACESDK_API __declspec(dllimport)
  #endif
#else
  #define CHLFACESDK_API	__attribute__((visibility("default")))
  #define WINAPI
#endif

#ifdef __cplusplus
  extern "C" {
#endif

// 获取SDK版本（可不调用）
// 输入参数：无
// 输出参数：无
// 返回值：SDK版本号
// 备注：可不调用，任何时候均可调用
CHLFACESDK_API int WINAPI ChlFaceSdkVer(void);

// SDK初始化
// 输入参数：
//     nMaxChannelNum ---- 需要支持的通道数，多线程应用每线程支持一个或多个通道
//     pLibPath ---- 依赖库文件所在目录
//     pTmpDir ---- 具有可写权限的临时目录
//     bOrigScore ---- 输出分数是否为原始分，为TRUE表示直接输出原始分(建议阈值55-65)，否则输出分数上浮15分(建议阈值70-80)
// 输出参数：无
// 返回值：成功返回0，许可无效返回-1，算法初始化失败返回-2
// 备注：除获取SDK版本外的任何接口都必须在SDK初始化成功后才能调用
//       ChlFaceSdkInit 输出分数上浮15分(建议阈值70-80)
CHLFACESDK_API int WINAPI ChlFaceSdkInit(int nMaxChannelNum, const char* pLibPath, const char* pTmpDir);
CHLFACESDK_API int WINAPI ChlFaceSdkInitEx(int nMaxChannelNum, const char* pLibPath, const char* pTmpDir, BOOL bOrigScore);


// SDK反初始化
// 输入参数：无
// 输出参数：无
// 返回值：无
// 备注：必须在初始化成功后调用，反初始化后不能再调用除获取SDK版本及SDK初始化外的任何接口
CHLFACESDK_API void WINAPI ChlFaceSdkUninit();

// 检测人脸
// 输入参数：
//           nChannelNo ---- 通道号（从0开始计算，限制在初始化指定的通道数内）
//           pRgb24 ---- RGB24格式的图象数据
//           nWidth ---- 图象数据宽度（象素单位）
//           nHeight ---- 图象数据高度（象素单位）
// 输出参数：pFaceResult ---- 检测到的人脸参数（人脸及眼睛等坐标位置及角度等，调用前必须分配有效的空间）
// 返回值：返回1表示检测到人脸，0表示无人脸，<0表示检测失败
CHLFACESDK_API int WINAPI ChlFaceSdkDetectFace(int nChannelNo, BYTE* pRgb24, int nWidth, int nHeight, LPFACE_DETECT_RESULT pFaceResult);

// 增强型检测人脸
// 输入参数：
//           nChannelNo ---- 通道号（从0开始计算，限制在初始化指定的通道数内）
//           pRgb24 ---- RGB24格式的图象数据
//           nWidth ---- 图象数据宽度（象素单位）
//           nHeight ---- 图象数据高度（象素单位）
//           nDetectSize ---- 检测图象大小（象素单位），如果图象宽度或高度大于此参数，则自动压缩到此参数指定大小以提供检测速度；此参数为0表示不压缩
// 输出参数：pFaceResult ---- 检测到的人脸参数（人脸及眼睛等坐标位置及角度等，调用前必须分配有效的空间）
// 返回值：返回1表示检测到人脸，0表示无人脸，<0表示检测失败
CHLFACESDK_API int WINAPI ChlFaceSdkDetectFaceEx(int nChannelNo, BYTE* pRgb24, int nWidth, int nHeight, int nDetectSize, LPFACE_DETECT_RESULT pFaceResult);

// 获取特征码大小
// 输入参数：无
// 输出参数：无
// 返回值：特征码大小
CHLFACESDK_API int WINAPI ChlFaceSdkFeatureSize(void);

// 提取特征码
// 输入参数：
//           nChannelNo ---- 通道号（从0开始计算，限制在初始化指定的通道数内）
//           pRgb24 ---- RGB24格式的图象数据
//           nWidth ---- 图象数据宽度
//           nHeight ---- 图象数据高度
//           pFaceResult ---- 检测到的人脸参数（必须将检测人脸返回的结果原样传入）
// 输出参数：pFeature ---- 特征码（调用前必须分配有效的空间）
// 返回值：成功返回0，失败返回-1
CHLFACESDK_API int WINAPI ChlFaceSdkFeatureGet(int nChannelNo, BYTE* pRgb24, int nWidth, int nHeight, LPFACE_DETECT_RESULT pFaceResult, BYTE* pFeature);

// 一对一特征比对
// 输入参数：
//           nChannelNo ---- 通道号（从0开始计算，限制在初始化指定的通道数内）
//           pFeature1 ---- 第1个人脸特征
//           pFeature2 ---- 第2个人脸特征
// 输出参数：无
// 返回值：返回两个人脸特征对应的人脸的相似度（0-100）
CHLFACESDK_API BYTE WINAPI ChlFaceSdkFeatureCompare(int nChannelNo, BYTE* pFeature1, BYTE* pFeature2);

// 检测多人脸并提特征(支持一张照片有多个人脸，多个人脸时按人脸大小排序，最大的人脸在最前面)
// 输入参数：
//           nChannelNo ---- 通道号（从0开始计算，限制在初始化指定的通道数内）
//           pRgb24 ---- RGB24格式的图象数据
//           nWidth ---- 图象数据宽度（象素单位）
//           nHeight ---- 图象数据高度（象素单位）
//           nMaxFace ---- 最多支持的人脸个数（1-10）
// 输出参数：pFaceResultBuf ---- 检测到的人脸参数（人脸及眼睛等坐标位置及角度等，调用前必须分配不小于 nMaxFace * sizeof(FACE_DETECT_RESULT) 的空间）
//           pFeaturesBuf ---- 检测到的人脸参数（人脸及眼睛等坐标位置及角度等，调用前必须分配不少于 nMaxFace * 特征大小 的空间）
// 返回值：返回检测到的人脸个数，0表示无人脸，-1表示检测失败
CHLFACESDK_API int WINAPI ChlFaceSdkFaceFeature(int nChannelNo, BYTE* pRgb24, int nWidth, int nHeight, int nMaxFace, LPFACE_DETECT_RESULT pFaceResultBuf, BYTE* pFeaturesBuf);


// 增强型检测多人脸并提特征(支持一张照片有多个人脸，多个人脸时按人脸大小排序，最大的人脸在最前面)
// 输入参数：
//           nChannelNo ---- 通道号（从0开始计算，限制在初始化指定的通道数内）
//           pRgb24 ---- RGB24格式的图象数据
//           nWidth ---- 图象数据宽度（象素单位）
//           nHeight ---- 图象数据高度（象素单位）
//           nDetectSize ---- 检测图象大小（象素单位），如果图象宽度或高度大于此参数，则自动压缩到此参数指定大小以提供检测速度；此参数为0表示不压缩
//           nMaxFace ---- 最多支持的人脸个数（1-10）
// 输出参数：pFaceResultBuf ---- 检测到的人脸参数（人脸及眼睛等坐标位置及角度等，调用前必须分配不小于 nMaxFace * sizeof(FACE_DETECT_RESULT) 的空间）
//           pFeaturesBuf ---- 检测到的人脸参数（人脸及眼睛等坐标位置及角度等，调用前必须分配不少于 nMaxFace * 特征大小 的空间）
// 返回值：返回检测到的人脸个数，0表示无人脸，-1表示检测失败
CHLFACESDK_API int WINAPI ChlFaceSdkFaceFeatureEx(int nChannelNo, BYTE* pRgb24, int nWidth, int nHeight, int nDetectSize, int nMaxFace, LPFACE_DETECT_RESULT pFaceResultBuf, BYTE* pFeaturesBuf);

// 检测人脸属性（性别年龄等, 支持同时对多人脸检测属性）
// 输入参数：
//           nChannelNo ---- 通道号（从0开始计算，限制在初始化指定的通道数内）
//           pRgb24 ---- RGB24格式的图象数据
//           nWidth ---- 图象数据宽度（象素单位）
//           nHeight ---- 图象数据高度（象素单位）
//           nFaceNum ---- 人脸个数
//           pFaceResult ---- 每个人脸的参数（人脸及眼睛等坐标位置等）
// 输出参数：
//           pPropertyBuf ---- 提取的人脸属性（性别年龄等，调用前必须分配不少于 nMaxFace * sizeof(FACE_PROPERTY_RESULT) 的空间）
// 返回值：成功返回0，-1表示参数无效，-2表示检测失败
CHLFACESDK_API int WINAPI ChlFaceSdkFaceProperty(int nChannelNo, BYTE* pRgb24, int nWidth, int nHeight, int nFaceNum, LPFACE_DETECT_RESULT pFaceResult, LPFACE_PROPERTY_RESULT pPropertyBuf);

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                               //
//  以下为一对多比对接口，对于目标人员较多且相对固定的场景，效率远比循环调用一对一接口时要高得多                 //
//                                                                                                               //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 创建一对多特征比对列表
// 输入参数：nMaxFeatureNum ---- 列表能容纳的最大特征数
// 输出参数：无
// 返回值：返回特征比对列表句柄，为空表示创建列表失败
// 备注：能容纳的最大特征数越大，则消耗的内存越多，另外实际支持的最大特征数可能受加密狗限制
CHLFACESDK_API HANDLE WINAPI ChlFaceSdkListCreate(int nMaxFeatureNum);

// 向一对多特征比对列表加入目标特征
// 输入参数：
//        hList ---- 目标特征要加入的特征比对列表句柄
//        pnPos ---- 指针非空时存放特征插入位置，位置值为0表示插到列表最前面，为-1或大于当前列表中的特征数则表示插到列表最后面
//        pFeatures ---- 要插入的目标特征码，多个特征按特征码大小（IdFaceSdkFeatureSize函数接口返回的值）对齐
//        nFeatureNum ---- 要插入的特征数量
// 输出参数：pnPos ---- 指针非空时返回第一个特征的插入位置
// 返回值：返回当前列表的总特征数
CHLFACESDK_API int WINAPI ChlFaceSdkListInsert(HANDLE hList, int* pnPos, BYTE* pFeatures, int nFeatureNum);

// 从一对多特征比对列表中删除部分特征
// 输入参数：
//        hList ---- 要删除特征的特征比对列表句柄
//        nPos ---- 要删除的特征的起始位置
//        nFeatureNum ---- 要删除的特征数量
// 输出参数：无
// 返回值：返回当前列表的总特征数
CHLFACESDK_API int WINAPI ChlFaceSdkListRemove(HANDLE hList, int nPos, int nFeatureNum);

// 清空一对多特征比对列表中的所有特征
// 输入参数：
//        hList ---- 要清空的一对多特征比对列表句柄
// 输出参数：无
// 返回值：无
// 备注：调用后列表中的实际特征数量为0（列表能容纳的最大特征数保持不变）
CHLFACESDK_API void WINAPI ChlFaceSdkListClearAll(HANDLE hList);

// 一对多特征比对
// 输入参数：
//        nChannelNo ---- 通道号（从0开始计算，限制在初始化指定的通道数内）
//        hList ---- 存放要参与比较的目标特征库的特征比对列表句柄
//        pFeature ---- 要参与特征比对的源特征码
//        nPosBegin ---- 要参与比对的目标特征的起始位置（如果比对全部目标特征，则起始位置填0）
//        nCompareNum ---- 要参与比对的目标特征的数量（如果比对全部目标特征，则比对数量填0或-1或实际特征数或最大特征数）
// 输出参数：pnScores ---- 顺序存放与各目标特征进行比对的相似度（每个特征比对的结果均为一字节，值范围为0-100）
// 返回值：返回实际参与比对的特征数量，也是返回的相似度个数
CHLFACESDK_API int WINAPI ChlFaceSdkListCompare(int nChannelNo, HANDLE hList, BYTE* pFeature, int nPosBegin, int nCompareNum, BYTE* pnScores);

// 销毁一对多特征比对列表
// 输入参数：
//        hList ---- 要销毁的一对多特征比对列表句柄
// 输出参数：无
// 返回值：无
// 备注：特征比对列表不再使用时需调用此接口释放资源
CHLFACESDK_API void WINAPI ChlFaceSdkListDestroy(HANDLE hList);


///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                               //
//  以下为扩展接口                                                                                               //
//                                                                                                               //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////


// 读图象文件，支持BMP、JPEG、PNG文件
// 输入参数：
//        pFileName ---- 图象文件名
//        nBufSize ---- 接收图象数据的缓冲区大小，当接收缓冲区大小为0或者接收缓冲区指针为空时，只输出图象的宽度
//        nDepth ---- 要接收的图象数据的位深度（需要灰度图象时填8，需要RGB24图象时填24）
// 输出参数：
//        pOutDataBuf ---- 接收图象数据，必须分配足够大的缓冲区并且将分配的缓冲区大小通过nBufSize参数传入
//        pnWidth ---- 接收图象宽度
//        pnHeight ---- 接收图象高度
// 返回值：
//        0 ---- 成功
//        -1 ---- 参数错误
//        -2 ---- 打开文件失败
//        其它 ---- 图象解析失败
// 备注：不知道图象分辨率因为不确认应分配的缓冲区大小时，可以先填入空缓冲区指针检测图象分辨率，再分配合适的缓冲区重新调用获得图象数据
CHLFACESDK_API int WINAPI ReadImageFile(const char* pFileName, BYTE* pRgbBuffer, int nBufSize, int* pnWidth, int* pnHeight, int nDepth);

CHLFACESDK_API int WINAPI SaveJpegFile(const char* pFileName, BYTE* pRgbBuffer, int nWidth, int nHeight, int nDepth, int nQuality);

#ifdef __cplusplus
  }
#endif

#endif // __CHLFACESDK__

