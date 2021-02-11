import { useRef } from 'react';
import { useIntersection } from 'react-use';
import gsap from 'gsap';

export type DirectionType = 'bottom' | 'left' | 'right' | 'top' | 'scale';

const fadeIn = (element: gsap.TweenTarget) => {
  gsap.to(element, 1, {
    opacity: 1,
    x: 0,
    y: 0,
    scale: 1,
    ease: 'power4.out',
  });
};

const fadeOut = (element: gsap.TweenTarget, direction: DirectionType) => {
  let x = 0;
  let y = 0;
  let scale = 1;
  switch (direction) {
    case 'bottom':
      y = +100;
      break;
    case 'left':
      x = +100;
      break;
    case 'right':
      x = -100;
      break;
    case 'scale':
      scale = 0.8;
      break;
    default:
      y = +100;
  }

  gsap.to(element, 1, {
    opacity: 0,
    x,
    y,
    scale,
    ease: 'power4.out',
  });
};

const useAnimation = (direction: DirectionType, intersectionParams: { rootMargin?: string }) => {
  const ref = useRef(null);
  const intersection = useIntersection(ref, {
    root: null,
    ...intersectionParams,
  });

  if (
    intersection &&
    ref.current &&
    (direction !== 'bottom' || intersection.boundingClientRect.y > 0) // avoid intersection box issue
  ) {
    intersection.isIntersecting ? fadeIn(ref.current) : fadeOut(ref.current, direction);
  }
  return ref;
};

export default useAnimation;
